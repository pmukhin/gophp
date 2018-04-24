package token

type (
	TokenType uint8
	Token struct {
		Type    TokenType
		Literal string
		Pos     int
	}
)

const (
	EOF                      TokenType = iota
	ILLEGAL
	NUMBER
	STRING
	SEMICOLON
	NOT
	IDENT
	EQUAL
	BACKSLASH
	COMMA
	COLON
	INCLUDE                   /* "include"			*/
	INCLUDE_ONCE              /* "include_once"			*/
	CURLY_OPENING
	CURLY_CLOSING
	REQUIRE                   /* "require"			*/
	REQUIRE_ONCE              /* "require_once"			*/
	LOGICAL_OR                /* "or"			*/
	LOGICAL_XOR               /* "xor"			*/
	LOGICAL_AND               /* "and"			*/
	SQUARE_BRACKET_OPENING
	SQUARE_BRACKET_CLOSING
	PRINT                     /* "print"			*/
	YIELD                     /* "yield"			*/
	YIELD_FROM                /* "yield from"			*/
	PLUS                      /* + */
	MINUS                     /* - */
	PLUS_EQUAL                /* "+="			*/
	MINUS_EQUAL               /* "-="			*/
	MUL
	MUL_EQUAL                 /* "*="			*/
	DIV
	DIV_EQUAL                 /* "/="			*/
	CONCAT_EQUAL              /* ".="			*/
	MOD
	MOD_EQUAL                 /* "%="			*/
	AND_EQUAL                 /* "&="			*/
	OR_EQUAL                  /* "|="			*/
	XOR_EQUAL                 /* "^="			*/
	SL_EQUAL                  /* "<<="			*/
	SR_EQUAL                  /* ">>="			*/
	BOOLEAN_OR                /* "||"			*/
	BOOLEAN_AND               /* "&&"			*/
	IS_EQUAL                  /* "=="			*/
	IS_NOT_EQUAL              /* "!="			*/
	IS_IDENTICAL              /* "==="			*/
	IS_NOT_IDENTICAL          /* "!=="			*/
	IS_SMALLER
	IS_GREATER
	IS_SMALLER_OR_EQUAL       /* "<="			*/
	IS_GREATER_OR_EQUAL       /* ">="			*/
	SPACESHIP                 /* "<=>"			*/
	SL                        /* "<<"			*/
	SR                        /* ">>"			*/
	INSTANCEOF                /* "instanceof"			*/
	INC                       /* "++"			*/
	DEC                       /* "--"			*/
	PARENTHESIS_OPENING
	PARENTHESIS_CLOSING
	INT_CAST                  /* "(int)"			*/
	DOUBLE_CAST               /* "(double)"			*/
	STRING_CAST               /* "(string)"			*/
	ARRAY_CAST                /* "(array)"			*/
	OBJECT_CAST               /* "(object)"			*/
	BOOL_CAST                 /* "(bool)"			*/
	UNSET_CAST                /* "(unset)"			*/
	NEW                       /* "new"			*/
	CLONE                     /* "clone"			*/
	EXIT                      /* "exit"			*/
	IF                        /* "if"			*/
	ELSEIF                    /* "elseif"			*/
	ELSE                      /* "else"			*/
	ENDIF                     /* "endif"			*/
	ECHO                      /* "echo"			*/
	DO                        /* "do"			*/
	WHILE                     /* "while"			*/
	ENDWHILE                  /* "endwhile"			*/
	FOR                       /* "for"			*/
	FOREACH                   /* "foreach"			*/
	DECLARE                   /* "declare"			*/
	ENDDECLARE                /* "enddeclare"			*/
	AS                        /* "as"			*/
	SWITCH                    /* "switch"			*/
	ENDSWITCH                 /* "endswitch"			*/
	CASE                      /* "case"			*/
	DEFAULT                   /* "default"			*/
	BREAK                     /* "break"			*/
	CONTINUE                  /* "continue"			*/
	GOTO                      /* "goto"			*/
	FUNCTION                  /* "function"			*/
	CONST                     /* "const"			*/
	RETURN                    /* "return"			*/
	TRY                       /* "try"			*/
	CATCH                     /* "catch"			*/
	FINALLY                   /* "finally"			*/
	THROW                     /* "throw"			*/
	USE                       /* "use"			*/
	INSTEADOF                 /* "insteadof"			*/
	GLOBAL                    /* "global"			*/
	STATIC                    /* "static"			*/
	ABSTRACT                  /* "abstract"			*/
	FINAL                     /* "final"			*/
	PRIVATE                   /* "private"			*/
	PROTECTED                 /* "protected"			*/
	PUBLIC                    /* "public"			*/
	VAR
	UNSET                     /* "unset"			*/
	ISSET                     /* "isset"			*/
	EMPTY                     /* "empty"			*/
	HALT_COMPILER             /* "__halt_compiler"			*/
	CLASS                     /* "class"			*/
	TRAIT                     /* "trait"			*/
	INTERFACE                 /* "interface"			*/
	EXTENDS                   /* "extends"			*/
	IMPLEMENTS                /* "implements"			*/
	OBJECT_OPERATOR           /* "->"			*/
	DOUBLE_ARROW              /* "=>"			*/
	DOUBLE_DOT                /* .. */
	LIST                      /* "list"			*/
	ARRAY                     /* "array"			*/
	CALLABLE                  /* "callable"			*/
	LINE                      /* "__LINE__"			*/
	FILE                      /* "__FILE__"			*/
	DIR                       /* "__DIR__"			*/
	CLASS_C                   /* "__CLASS__"			*/
	TRAIT_C                   /* "__TRAIT__"			*/
	METHOD_C                  /* "__METHOD__"			*/
	FUNC_C                    /* "__FUNCTION__"			*/
	COMMENT_START
	COMMENT_END
	COMMENT                   /* "comment"			*/
	DOC_COMMENT               /* "doc comment"			*/
	OPEN_TAG                  /* "open tag"			*/
	OPEN_TAG_WITH_ECHO        /* "open tag with echo"			*/
	CLOSE_TAG                 /* "close tag"			*/
	WHITESPACE                /* "whitespace"			*/
	START_HEREDOC             /* "heredoc start"			*/
	END_HEREDOC               /* "heredoc end"			*/
	DOLLAR_OPEN_CURLY_BRACES  /* "${"			*/
	CURLY_OPEN                /* "{$"			*/
	PAAMAYIM_NEKUDOTAYIM      /* "::"			*/
	NAMESPACE                 /* "namespace"			*/
	NS_C                      /* "__NAMESPACE__"			*/
	NS_SEPARATOR              /* "\\"			*/
	ELLIPSIS                  /* "..."			*/
	COALESCE                  /* "??"			*/
	POW                       /* "**"			*/
	POW_EQUAL                 /* "**="			*/
	NEWLINE
)
