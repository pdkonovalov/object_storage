package object_storage

type Object struct {
	Path     string
	Meta     map[string]any
	Contains []string
	URL      string
}
