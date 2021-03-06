package application

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	log_helper "github.com/krafton-hq/golib/log-helper"
	"github.com/krafton-hq/red-fox/apis/api_resources"
	"github.com/krafton-hq/red-fox/apis/app_lifecycle"
	"github.com/krafton-hq/red-fox/apis/crds"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/krafton-hq/red-fox/apis/namespaces"
	"github.com/krafton-hq/red-fox/server/application/configs"
	"github.com/krafton-hq/red-fox/server/controllers/api_resources_con"
	"github.com/krafton-hq/red-fox/server/controllers/app_lifecycle_con"
	"github.com/krafton-hq/red-fox/server/controllers/crd_con"
	"github.com/krafton-hq/red-fox/server/controllers/document_con"
	"github.com/krafton-hq/red-fox/server/controllers/external_dns_con"
	"github.com/krafton-hq/red-fox/server/controllers/namespace_con"
	"github.com/krafton-hq/red-fox/server/pkg/database_helper"
	"github.com/krafton-hq/red-fox/server/pkg/domain_helper"
	"github.com/krafton-hq/red-fox/server/pkg/errors"
	"github.com/krafton-hq/red-fox/server/pkg/transactional"
	"github.com/krafton-hq/red-fox/server/repositories/apiobject_repository"
	"github.com/krafton-hq/red-fox/server/repositories/repository_manager"
	"github.com/krafton-hq/red-fox/server/services/crd_service"
	"github.com/krafton-hq/red-fox/server/services/custom_document_service"
	"github.com/krafton-hq/red-fox/server/services/endpoint_service"
	"github.com/krafton-hq/red-fox/server/services/external_dns_service"
	"github.com/krafton-hq/red-fox/server/services/namespace_service"
	"github.com/krafton-hq/red-fox/server/services/natip_service"
	"github.com/krafton-hq/red-fox/server/services/service_decorator"
	"github.com/krafton-hq/red-fox/server/services/services"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

const Version = "0.0.0-dev-build"

type Application struct {
	config *configs.RedFoxConfig

	grpcServer *grpc.Server

	nsController           *namespace_con.Controller
	crdController          *crd_con.Controller
	natIpController        *document_con.NatIpController
	endpointController     *document_con.EndpointController
	customDocController    *document_con.CustomDocumentController
	apiResourcesController *api_resources_con.Controller
	appController          *app_lifecycle_con.GrpcController
	extDnsController       *external_dns_con.Controller

	extDnsService *external_dns_service.Service
}

func NewApplication(config *configs.RedFoxConfig) *Application {
	return &Application{config: config}
}

func (a *Application) Init() error {
	err := a.initInternal()
	if err != nil {
		return err
	}

	var panicHandler grpc_recovery.RecoveryHandlerFuncContext = func(ctx context.Context, p interface{}) (err error) {
		zap.S().Warnw("Panic Detected", "error", p, "stacktrace", string(debug.Stack()))
		return status.Errorf(codes.Internal, "%v", p)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			append(log_helper.GetUnaryServerInterceptors(),
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(panicHandler)),
			)...,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			append(log_helper.GetStreamServerInterceptors(),
				grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(panicHandler)),
			)...,
		)))
	reflection.Register(grpcServer)

	app_lifecycle.RegisterApplicationLifecycleServer(grpcServer, a.appController)
	namespaces.RegisterNamespaceServerServer(grpcServer, a.nsController)
	crds.RegisterCustomResourceDefinitionServerServer(grpcServer, a.crdController)
	documents.RegisterNatIpServerServer(grpcServer, a.natIpController)
	documents.RegisterEndpointServerServer(grpcServer, a.endpointController)
	documents.RegisterCustomDocumentServerServer(grpcServer, a.customDocController)
	api_resources.RegisterApiResourcesServerServer(grpcServer, a.apiResourcesController)

	for name := range grpcServer.GetServiceInfo() {
		zap.S().Infow("Registered gRpc Service", "name", name)
	}

	wrappedGrpc := grpcweb.WrapServer(grpcServer)
	grpcWebHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			if strings.HasPrefix(req.URL.Path, "/grpc_web") {
				req.URL.Path = strings.Replace(req.URL.Path, "/grpc_web", "", 1)
			}
			wrappedGrpc.ServeHTTP(res, req)
			return
		}
		zap.S().Info("Unknown request received")
		res.WriteHeader(http.StatusNotFound)
	})
	httpServer := fiber.New()
	httpServer.Use(logger.New())
	httpServer.Group("/grpc_web", adaptor.HTTPHandlerFunc(grpcWebHandler))

	appLifecycle := httpServer.Group("/api/v1/app")
	app_lifecycle_con.NewAppLifecycleHttp(a.appController).Register(appLifecycle)

	go func(port int32) {
		listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		fmt.Printf("Grpc Server listen http://0.0.0.0:%d\n", port)
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}(a.config.Listeners.GrpcPort)

	go func(port int32) {
		fmt.Printf("Grpc Web Server listen http://0.0.0.0:%d\n", port)
		err := httpServer.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to start grpc-web server: %v", err)
		}
	}(a.config.Listeners.RestPort)

	if a.config.ExternalDns.Enabled {
		go a.extDnsService.Run(context.TODO())

		go func(port int32) {
			fmt.Printf("External-Dns TCP Server listen tcp://0.0.0.0:%d\n", port)
			ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				log.Fatalf("failed to Listen External-Dns TCP server: %v", err)
			}
			err = a.extDnsController.Start(ln)
			if err != nil {
				log.Fatalf("failed to Start External-Dns TCP server: %v", err)
			}
		}(a.config.ExternalDns.Port)
	}

	return nil
}

func (a *Application) initInternal() error {

	var tr transactional.Transactional
	switch a.config.Database.Type {
	case configs.DatabaseType_Inmemory:
		tr = transactional.NewNoop()
	case configs.DatabaseType_Mysql:
		db, err := database_helper.NewDatabase(a.config.Database.Url, configs.ParseStringRef(a.config.Database.UsernameRef), configs.ParseStringRef(a.config.Database.PasswordRef))
		if err != nil {
			zap.S().Errorw("Create Database Connection Failed", "error", err)
			return err
		}
		tr, err = transactional.NewSqlTransactional(db, transactional.DialectMysql, nil)
		if err != nil {
			zap.S().Errorw("Create Database Transaction Helper Failed", "error", err)
			return err
		}
	default:
		return errors.NewErrorf("Unknown Database Type: %s", a.config.Database.Type.String())
	}

	var nsRepo apiobject_repository.ClusterRepository[*namespaces.Namespace]
	var crdRepo apiobject_repository.ClusterRepository[*crds.CustomResourceDefinition]
	var natIpRepo apiobject_repository.NamespacedRepository[*documents.NatIp]
	var endpointRepo apiobject_repository.NamespacedRepository[*documents.Endpoint]
	var customDocRepoFactory apiobject_repository.ClusterRepositoryFactory[*documents.CustomDocument]

	switch a.config.Database.Type {
	case configs.DatabaseType_Inmemory:
		nsRepo = apiobject_repository.NewInMemoryClusterRepository[*namespaces.Namespace](domain_helper.NamespaceGvk, apiobject_repository.DefaultSystemNamespace)
		crdRepo = apiobject_repository.NewInMemoryClusterRepository[*crds.CustomResourceDefinition](domain_helper.CrdGvk, apiobject_repository.DefaultSystemNamespace)
		natIpRepo = apiobject_repository.NewGenericNamespacedRepository[*documents.NatIp](domain_helper.NatIpGvk, apiobject_repository.NewInmemoryClusterRepositoryFactory[*documents.NatIp]())
		endpointRepo = apiobject_repository.NewGenericNamespacedRepository[*documents.Endpoint](domain_helper.EndpointGvk, apiobject_repository.NewInmemoryClusterRepositoryFactory[*documents.Endpoint]())
		customDocRepoFactory = apiobject_repository.NewInmemoryClusterRepositoryFactory[*documents.CustomDocument]()
		break
	case configs.DatabaseType_Mysql:
		nsRepo = apiobject_repository.NewMysqlClusterRepository[*namespaces.Namespace](domain_helper.NamespaceGvk, apiobject_repository.DefaultSystemNamespace, domain_helper.NewNamespaceFactory(), tr)
		crdRepo = apiobject_repository.NewMysqlClusterRepository[*crds.CustomResourceDefinition](domain_helper.CrdGvk, apiobject_repository.DefaultSystemNamespace, domain_helper.NewCrdFactory(), tr)
		natIpRepo = apiobject_repository.NewGenericNamespacedRepository[*documents.NatIp](domain_helper.NatIpGvk, apiobject_repository.NewMysqlClusterRepositoryFactory[*documents.NatIp](domain_helper.NewNatIpFactory(), tr))
		endpointRepo = apiobject_repository.NewGenericNamespacedRepository[*documents.Endpoint](domain_helper.EndpointGvk, apiobject_repository.NewMysqlClusterRepositoryFactory[*documents.Endpoint](domain_helper.NewEndpointFactory(), tr))
		customDocRepoFactory = apiobject_repository.NewMysqlClusterRepositoryFactory[*documents.CustomDocument](domain_helper.NewCustomDocumentFactory(), tr)
		break
	default:
		return errors.NewErrorf("Unknown Database Type: %s", a.config.Database.Type.String())
	}

	var natIpService services.NamespacedService[*documents.NatIp] = natip_service.NewService(natIpRepo)
	var endpointService services.NamespacedService[*documents.Endpoint] = endpoint_service.NewService(endpointRepo)
	repoManager := repository_manager.NewManager(customDocRepoFactory, natIpRepo, endpointRepo, nsRepo, crdRepo)
	var customDocService services.NamespacedGvkService[*documents.CustomDocument] = custom_document_service.NewService(repoManager)
	var nsService services.ClusterService[*namespaces.Namespace] = namespace_service.NewService(nsRepo, repoManager)
	var crdService services.ClusterService[*crds.CustomResourceDefinition] = crd_service.NewService(crdRepo, repoManager)

	nsService = service_decorator.NewTransactionalClusterService[*namespaces.Namespace](nsService, tr)
	crdService = service_decorator.NewTransactionalClusterService[*crds.CustomResourceDefinition](crdService, tr)
	natIpService = service_decorator.NewTransactionalNamespacedService[*documents.NatIp](natIpService, tr)
	endpointService = service_decorator.NewTransactionalNamespacedService[*documents.Endpoint](endpointService, tr)
	customDocService = service_decorator.NewTransactionalNamespacedGvkService[*documents.CustomDocument](customDocService, tr)

	a.nsController = namespace_con.NewController(nsService)
	a.crdController = crd_con.NewController(crdService)
	a.natIpController = document_con.NewNatIpDocController(natIpService)
	a.endpointController = document_con.NewEndpointController(endpointService)
	a.customDocController = document_con.NewCustomDocumentController(customDocService)
	a.appController = app_lifecycle_con.NewAppLifecycle()
	a.apiResourcesController = api_resources_con.NewController(repoManager)

	if a.config.ExternalDns.Enabled {
		duration, err := time.ParseDuration(a.config.ExternalDns.SyncInterval)
		if err != nil {
			return err
		}

		tmplCfg := a.config.ExternalDns.Templates
		a.extDnsService = external_dns_service.NewService(&external_dns_service.Templates{
			NatIpName:              tmplCfg.NatIpName,
			NatIpLabel:             tmplCfg.NatIpLabel,
			NatIpLabelWithValue:    tmplCfg.NatIpLabelWithValue,
			EndpointName:           tmplCfg.EndpointName,
			EndpointLabel:          tmplCfg.EndpointLabel,
			EndpointLabelWithValue: tmplCfg.EndpointLabelWithValue,
		}, duration, natIpService)
		a.extDnsController = external_dns_con.NewController(a.extDnsService)
	}

	err := crdService.Init(context.TODO())
	if err != nil {
		return err
	}
	err = nsService.Init(context.TODO())
	if err != nil {
		return err
	}
	return nil
}
