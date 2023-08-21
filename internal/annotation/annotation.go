package annotation

const (
	ManagedByAnnotation = "octopipe.io/managed-by"
)

var DefaultAnnotations = map[string]string{
	ManagedByAnnotation: "cloudx",
}
