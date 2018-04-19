package ast

import (
	"testing"
)

func TestBlockStatement_String(t *testing.T) {
	type fields struct {
		Statements []Statement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BlockStatement{
				Statements: tt.fields.Statements,
			}
			if got := b.String(); got != tt.want {
				t.Errorf("BlockStatement.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceOfExpression_String(t *testing.T) {
	type fields struct {
		Object Expression
		Type   Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ioe := InstanceOfExpression{
				Object: tt.fields.Object,
				Type:   tt.fields.Type,
			}
			if got := ioe.String(); got != tt.want {
				t.Errorf("InstanceOfExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionalExpression_String(t *testing.T) {
	type fields struct {
		Condition   Expression
		Consequence *BlockStatement
		Alternative *BlockStatement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConditionalExpression{
				Condition:   tt.fields.Condition,
				Consequence: tt.fields.Consequence,
				Alternative: tt.fields.Alternative,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("ConditionalExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssignmentExpression_String(t *testing.T) {
	type fields struct {
		Left  Expression
		Right Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ae := AssignmentExpression{
				Left:  tt.fields.Left,
				Right: tt.fields.Right,
			}
			if got := ae.String(); got != tt.want {
				t.Errorf("AssignmentExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNull_String(t *testing.T) {
	tests := []struct {
		name string
		n    Null
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Null{}
			if got := n.String(); got != tt.want {
				t.Errorf("Null.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariableExpression_String(t *testing.T) {
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ve := VariableExpression{
				Name: tt.fields.Name,
			}
			if got := ve.String(); got != tt.want {
				t.Errorf("VariableExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegerLiteral_String(t *testing.T) {
	type fields struct {
		Value int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			il := IntegerLiteral{
				Value: tt.fields.Value,
			}
			if got := il.String(); got != tt.want {
				t.Errorf("IntegerLiteral.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringLiteral_Pos(t *testing.T) {
	type fields struct {
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringLiteral{
				Value: tt.fields.Value,
			}
			if got := s.Pos(); got != tt.want {
				t.Errorf("StringLiteral.Pos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringLiteral_End(t *testing.T) {
	type fields struct {
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringLiteral{
				Value: tt.fields.Value,
			}
			if got := s.End(); got != tt.want {
				t.Errorf("StringLiteral.End() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringLiteral_TokenLiteral(t *testing.T) {
	type fields struct {
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringLiteral{
				Value: tt.fields.Value,
			}
			if got := s.TokenLiteral(); got != tt.want {
				t.Errorf("StringLiteral.TokenLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringLiteral_String(t *testing.T) {
	type fields struct {
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := StringLiteral{
				Value: tt.fields.Value,
			}
			if got := sl.String(); got != tt.want {
				t.Errorf("StringLiteral.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringLiteral_expressionNode(t *testing.T) {
	type fields struct {
		Value string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringLiteral{
				Value: tt.fields.Value,
			}
			s.expressionNode()
		})
	}
}

func TestArrayLiteral_String(t *testing.T) {
	type fields struct {
		Elements []Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			al := ArrayLiteral{
				Elements: tt.fields.Elements,
			}
			if got := al.String(); got != tt.want {
				t.Errorf("ArrayLiteral.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexExpression_String(t *testing.T) {
	type fields struct {
		Index Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple integer",
			fields: fields{
				Index: IntegerLiteral{Value: 0},
			},
			want: "[0]",
		},
		{
			name: "simple string",
			fields: fields{
				Index: StringLiteral{Value: "test"},
			},
			want: "['test']",
		},
		{
			name: "arithmetic expression",
			fields: fields{
				Index: BinaryExpression{
					Left:  IntegerLiteral{Value: 0},
					Right: IntegerLiteral{Value: 1},
					Op:    "+",
				},
			},
			want: "[0 + 1]",
		},
		{
			name: "complex expression",
			fields: fields{
				Index: BinaryExpression{
					Left:  VariableExpression{Name: "integerVar"},
					Right: IntegerLiteral{Value: 1},
					Op:    "-",
				},
			},
			want: "[$integerVar - 1]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ie := IndexExpression{
				Index: tt.fields.Index,
			}
			if got := ie.String(); got != tt.want {
				t.Errorf("IndexExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooleanExpression_String(t *testing.T) {
	type fields struct {
		Value bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := BooleanExpression{
				Value: tt.fields.Value,
			}
			if got := be.String(); got != tt.want {
				t.Errorf("BooleanExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnaryExpression_String(t *testing.T) {
	type fields struct {
		IsPrefix bool
		Op       string
		Right    Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := UnaryExpression{
				IsPrefix: tt.fields.IsPrefix,
				Op:       tt.fields.Op,
				Right:    tt.fields.Right,
			}
			if got := pe.String(); got != tt.want {
				t.Errorf("UnaryExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryExpression_String(t *testing.T) {
	type fields struct {
		Left  Expression
		Op    string
		Right Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "integers",
			fields: fields{
				Left:  IntegerLiteral{Value: 5},
				Right: IntegerLiteral{Value: 5},
				Op:    "+",
			},
			want: "5 + 5",
		},
		{
			name: "strings",
			fields: fields{
				Left:  StringLiteral{Value: "Test1"},
				Right: StringLiteral{Value: "Test2"},
				Op:    "+",
			},
			want: "'Test1' + 'Test2'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			be := BinaryExpression{
				Left:  tt.fields.Left,
				Op:    tt.fields.Op,
				Right: tt.fields.Right,
			}
			if got := be.String(); got != tt.want {
				t.Errorf("BinaryExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhileExpression_String(t *testing.T) {
	type fields struct {
		Condition Expression
		Body      *BlockStatement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			we := WhileExpression{
				Condition: tt.fields.Condition,
				Body:      tt.fields.Body,
			}
			if got := we.String(); got != tt.want {
				t.Errorf("WhileExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForEachExpression_String(t *testing.T) {
	type fields struct {
		Array Expression
		Key   *VariableExpression
		Value *VariableExpression
		Block *BlockStatement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "with key",
			fields: fields{
				Array: ArrayLiteral{
					Elements: []Expression{
						IntegerLiteral{1},
						IntegerLiteral{2},
						IntegerLiteral{3},
					},
				},
				Key:   &VariableExpression{"key"},
				Value: &VariableExpression{"value"},
				Block: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: &AssignmentExpression{
								Left:  VariableExpression{Name: "one"},
								Right: VariableExpression{Name: "value"},
							},
						},
					},
				},
			},
			want: `foreach([1, 2, 3] as $key => $value) {
    $one = $value;
}`,
		},
		{
			name: "without key",
			fields: fields{
				Array: &VariableExpression{"array"},
				Value: &VariableExpression{"value"},
				Block: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{
							Expression: &AssignmentExpression{
								Left:  VariableExpression{Name: "one"},
								Right: VariableExpression{Name: "value"},
							},
						},
					},
				},
			},
			want: `foreach($array as $value) {
    $one = $value;
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee := ForEachExpression{
				Array: tt.fields.Array,
				Key:   tt.fields.Key,
				Value: tt.fields.Value,
				Block: tt.fields.Block,
			}
			if got := fee.String(); got != tt.want {
				t.Errorf("ForEachExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpressionStatement_String(t *testing.T) {
	type fields struct {
		Expression Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := ExpressionStatement{
				Expression: tt.fields.Expression,
			}
			if got := es.String(); got != tt.want {
				t.Errorf("ExpressionStatement.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionCall_String(t *testing.T) {
	type fields struct {
		Target   Identifier
		CallArgs []Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty args",
			fields: fields{
				Target: Identifier{Value: "php_info"},
			},
			want: "php_info()",
		},
		{
			name: "with args",
			fields: fields{
				Target: Identifier{Value: "func"},
				CallArgs: []Expression{
					FunctionCall{Target: Identifier{Value: "some"}, CallArgs: []Expression{IntegerLiteral{Value: 1}}},
					StringLiteral{Value: "TestString"},
					BooleanExpression{Value: false},
				},
			},
			want: "func(some(1), 'TestString', false)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := FunctionCall{
				Target:   tt.fields.Target,
				CallArgs: tt.fields.CallArgs,
			}
			if got := fc.String(); got != tt.want {
				t.Errorf("FunctionCall.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMethodCall_String(t *testing.T) {
	type fields struct {
		FunctionCall FunctionCall
		Object       Expression
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "with args",
			fields: fields{
				FunctionCall: FunctionCall{
					Target: Identifier{Value: "func"},
					CallArgs: []Expression{
						FunctionCall{Target: Identifier{Value: "some"}, CallArgs: []Expression{IntegerLiteral{Value: 1}}},
						StringLiteral{Value: "TestString"},
						BooleanExpression{Value: false},
					},
				},
				Object: VariableExpression{Name: "object"},
			},
			want: "$object->func(some(1), 'TestString', false)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := MethodCall{
				FunctionCall: tt.fields.FunctionCall,
				Object:       tt.fields.Object,
			}
			if got := mc.String(); got != tt.want {
				t.Errorf("MethodCall.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentifier_String(t *testing.T) {
	type fields struct {
		Value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Identifier{
				Value: tt.fields.Value,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("Identifier.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArg_String(t *testing.T) {
	type fields struct {
		Type         *Identifier
		Name         VariableExpression
		DefaultValue Expression
		Variadic     bool
		IsReference  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "typed arg reference",
			fields: fields{
				Type:        &Identifier{Value: "string"},
				Name:        VariableExpression{Name: "arg"},
				IsReference: true,
			},
			want: "string &$arg",
		},
		{
			name: "typed variadic arg",
			fields: fields{
				Type:     &Identifier{Value: "string"},
				Name:     VariableExpression{Name: "args"},
				Variadic: true,
			},
			want: "string ...$args",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Arg{
				Type:         tt.fields.Type,
				Name:         tt.fields.Name,
				DefaultValue: tt.fields.DefaultValue,
				Variadic:     tt.fields.Variadic,
				IsReference:  tt.fields.IsReference,
			}
			if got := a.String(); got != tt.want {
				t.Errorf("Arg.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFunctionDeclarationExpression_String(t *testing.T) {
	type fields struct {
		Name       Identifier
		Args       []Arg
		ReturnType *Identifier
		Block      *BlockStatement
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple function with args",
			fields: fields{
				Name:       Identifier{Value: "testFunction"},
				Args:       []Arg{{Type: &Identifier{Value: "int"}, Name: VariableExpression{Name: "var"}}},
				ReturnType: &Identifier{Value: "ReturnType"},
				Block: &BlockStatement{Statements: []Statement{
					ExpressionStatement{
						Expression: AssignmentExpression{Left: VariableExpression{Name: "var"}, Right: IntegerLiteral{Value: 1}},
					},
				}},
			},
			want: `function testFunction(int $var): ReturnType {
    $var = 1;
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fde := FunctionDeclarationExpression{
				Name:       tt.fields.Name,
				Args:       tt.fields.Args,
				ReturnType: tt.fields.ReturnType,
				Block:      tt.fields.Block,
			}
			if got := fde.String(); got != tt.want {
				t.Errorf("FunctionDeclarationExpression.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
