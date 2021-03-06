package object

// RegisterGlobals registers all built-in functions
func RegisterGlobals(ctx Context) error {
	registerPrintFunctions(ctx)
	registerMathFunctions(ctx)
	registerOsConstants(ctx)
	registerIntConstants(ctx)
	registerStringConstants(ctx)
	registerArrayConstants(ctx)

	return nil
}
