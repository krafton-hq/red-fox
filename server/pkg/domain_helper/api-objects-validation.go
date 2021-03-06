package domain_helper

import (
	"fmt"

	"github.com/krafton-hq/red-fox/server/pkg/errors"
	"github.com/krafton-hq/red-fox/server/pkg/validation"
)

const subdomainMaxLength int = 63
const nameMaxLength int = subdomainMaxLength
const labelKeyMaxLength int = subdomainMaxLength
const labelValueMaxLength int = 4095
const qualifiedNameMaxLength int = 253

const fieldName = "metadata.name"
const fieldLabelKey = "metadata.label[key]"
const fieldLabelValue = "metadata.label[value]"
const fieldAnnotation = "metadata.annotations"
const fieldAnnotationKey = "metadata.annotations[key]"

func ValidationMetadatable(m Metadatable) error {
	if m == nil {
		return errors.NewInvalidField("$", "Should not be null", "null")
	}

	metadata := m.GetMetadata()
	if errs := validation.IsApiVersion(m.GetApiVersion()); len(errs) > 0 {
		return errors.NewInvalidField("apiVersion", "RFC1123 Dns Label/Version", m.GetApiVersion())
	}

	for key, value := range metadata.Labels {
		if len(key) == 0 || len(key) > labelKeyMaxLength {
			return errors.NewInvalidField(fieldLabelKey, fmt.Sprintf("Label Key Length Should be [1, %d]", labelKeyMaxLength), key)
		}
		if errs := validation.IsDiscoveryName(key); len(errs) > 0 {
			return errors.NewInvalidField(fieldLabelKey, "RFC1123 Dns Label", key)
		}

		if len(value) > labelValueMaxLength {
			return errors.NewInvalidField(fieldLabelValue, fmt.Sprintf("Label Value Length Should be [0, %d]", labelValueMaxLength), value)
		}
	}

	for key := range metadata.Annotations {
		if errs := validation.IsAnnotationName(key); len(errs) > 0 {
			return errors.NewInvalidField(fieldAnnotationKey, fmt.Sprintf("%v", errs), key)
		}
	}
	if errs := validation.IsValidAnnotationsSize(metadata.Annotations); len(errs) > 0 {
		return errors.NewInvalidField(fieldAnnotation, fmt.Sprintf("%v", errs), fmt.Sprintf("%v", metadata.Annotations))
	}
	return nil
}

func ValidationQualifiedName(m Metadatable) error {
	if m == nil {
		return errors.NewInvalidField("$", "Should not be null", "null")
	}

	metadata := m.GetMetadata()
	if len(metadata.Name) == 0 || len(metadata.Name) > qualifiedNameMaxLength {
		return errors.NewInvalidField(fieldName, fmt.Sprintf("Length Should be [1, %d]", qualifiedNameMaxLength), metadata.Name)
	}
	if errs := validation.IsQualifiedName(metadata.Name); len(errs) > 0 {
		return errors.NewInvalidField(fieldName, "RFC1123 Domain Name", metadata.Name)
	}
	return nil
}

func ValidationDiscoverableName(m Metadatable) error {
	if m == nil {
		return errors.NewInvalidField("$", "Should not be null", "null")
	}

	metadata := m.GetMetadata()
	if len(metadata.Name) == 0 || len(metadata.Name) > nameMaxLength {
		return errors.NewInvalidField(fieldName, fmt.Sprintf("Length Should be [1, %d]", nameMaxLength), metadata.Name)
	}
	if errs := validation.IsDiscoveryName(metadata.Name); len(errs) > 0 {
		return errors.NewInvalidField(fieldName, "RFC1123 Dns Label", metadata.Name)
	}
	return nil
}
