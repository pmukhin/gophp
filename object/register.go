package object

// RegisterGlobals registers all built-in functions
func RegisterGlobals(ctx Context) error {
	registerPrintFunctions(ctx)
	registerMathFunctions(ctx)
	registerOsConstants(ctx)

	return nil
}
