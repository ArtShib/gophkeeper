package customprototype

import (
	keeperv1 "github.com/ArtShib/gophkeeper/gen/go/keeper/v1"
	"github.com/ArtShib/gophkeeper/internal/models"
)

func GetProtoSecretType(secretType models.SecretType) *keeperv1.SecretType {
	var protoType keeperv1.SecretType

	switch secretType {
	case models.TypeText:
		protoType = keeperv1.SecretType_SECRET_TYPE_TEXT
	case models.TypeBinary:
		protoType = keeperv1.SecretType_SECRET_TYPE_BINARY
	case models.TypeCard:
		protoType = keeperv1.SecretType_SECRET_TYPE_CARD
	case models.TypeCredentials:
		protoType = keeperv1.SecretType_SECRET_TYPE_CREDENTIALS
	default:
		protoType = keeperv1.SecretType_SECRET_TYPE_UNSPECIFIED
	}
	return &protoType
}

func ProtoToSecretType(protoType keeperv1.SecretType) models.SecretType {
	var secretType models.SecretType

	switch protoType {
	case keeperv1.SecretType_SECRET_TYPE_CREDENTIALS:
		secretType = models.TypeCredentials
	case keeperv1.SecretType_SECRET_TYPE_BINARY:
		secretType = models.TypeBinary
	case keeperv1.SecretType_SECRET_TYPE_CARD:
		secretType = models.TypeCard
	case keeperv1.SecretType_SECRET_TYPE_TEXT:
		secretType = models.TypeText
	default:
		secretType = models.TypeUnspecifed
	}
	return secretType
}
