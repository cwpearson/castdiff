// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cc

import (
	"fmt"
)

// An Expr is a parsed C expression.
type Expr struct {
	SyntaxInfo
	Id           int
	Op           ExprOp   // operator
	Left         *Expr    // left (or only) operand
	Right        *Expr    // right operand
	List         []*Expr  // operand list, for Comma, Cond, Call
	LaunchParams []*Expr  // launch params for CUDACall
	Text         Syntax   // name or literal, for Name, Number, Goto, Arrow, Dot
	Texts        []Syntax // list of literals, for String
	Type         *Type    // type operand, for SizeofType, Offsetof, Cast, CastInit, VaArg
	Init         *Init    // initializer, for CastInit
	Block        []*Stmt  // for c2go
	SourceExpr   *Expr
	// derived information
	XDecl *Decl
	XType *Type // expression type, derived
}

func (x *Expr) GetId() int {
	return x.Id
}

func (x *Expr) GetChildren() []Syntax {
	lst := []Syntax{}
	if x.Text != nil {
		lst = append(lst, x.Text)
	}
	if len(x.Texts) != 0 {
		lst = append(lst, x.Texts...)
	}
	switch x.Op {
	default:
		if x.Type != nil {
			lst = append(lst, x.Type)
		}
		if x.Left != nil {
			lst = append(lst, x.Left)
		}
		if x.Right != nil {
			lst = append(lst, x.Right)
		}
		if len(x.List) != 0 {
			for _, elem := range x.List {
				if elem != nil {
					lst = append(lst, elem)
				}
			}
		}
	case Arrow:
		lst = append(lst, x.Left)
	case Call:
		lst = append(lst, x.Left)
		if len(x.List) != 0 {
			for _, elem := range x.List {
				if elem != nil {
					lst = append(lst, elem)
				}
			}
		}
	case Comma:
		if len(x.List) != 0 {
			for _, elem := range x.List {
				if elem != nil {
					lst = append(lst, elem)
				}
			}
		}
	case Cond:
		lst = append(lst, x.List[0], x.List[1], x.List[2])
	case Dot:
		lst = append(lst, x.Left)
	case Cast:
		lst = append(lst, x.Type)
	case CastInit:
		lst = append(lst, x.Type)
	case Index:
		lst = append(lst, x.Left, x.Right)
	case Offsetof:
		lst = append(lst, x.Type, x.Left)
	case Paren:
		lst = append(lst, x.Left)
	case PostDec:
		lst = append(lst, x.Left)
	case PostInc:
		lst = append(lst, x.Left)
	case VaArg:
		lst = append(lst, x.Left, x.Type)
	}
	return lst
}

func (x *Expr) String() string {
	var p Printer
	p.hideComments = true
	p.printExpr(x, precLow)
	return p.String()
}

type ExprOp int

const (
	_          ExprOp = iota
	Add               // Left + Right
	AddEq             // Left += Right
	Addr              // &Left
	And               // Left & Right
	AndAnd            // Left && Right
	AndEq             // Left &= Right
	Arrow             // Left->Text
	Call              // Left(List)
	CUDACall          // Left(LaunchParams, List)
	Cast              // (Type)Left
	CastInit          // (Type){Init}
	Comma             // x, y, z; List = {x, y, z}
	Cond              // x ? y : z; List = {x, y, z}
	Div               // Left / Right
	DivEq             // Left /= Right
	Dot               // Left.Name
	Eq                // Left = Right
	EqEq              // Left == Right
	Gt                // Left > Right
	GtEq              // Left >= Right
	Index             // Left[Right]
	Indir             // *Left
	Lsh               // Left << Right
	LshEq             // Left <<= Right
	LcuBrk            // Left <<< Right
	Lt                // Left < Right
	LtEq              // Left <= Right
	RcuBrk            // Left >>> Right
	Minus             // -Left
	Mod               // Left % Right
	ModEq             // Left %= Right
	Mul               // Left * Right
	MulEq             // Left *= Right
	Name              // Text (function, variable, or enum name)
	Not               // !Left
	NotEq             // Left != Right
	Number            // Text (numeric or chraracter constant)
	Literal           // Text (numeric or chraracter constant)
	Offsetof          // offsetof(Type, Left)
	Or                // Left | Right
	OrEq              // Left |= Right
	OrOr              // Left || Right
	Paren             // (Left)
	Plus              //  +Left
	PostDec           // Left--
	PostInc           // Left++
	PreDec            // --Left
	PreInc            // ++Left
	Rsh               // Left >> Right
	RshEq             // Left >>= Right
	SizeofExpr        // sizeof(Left)
	SizeofType        // sizeof(Type)
	String            // Text (quoted string literal)
	Sub               // Left - Right
	SubEq             // Left -= Right
	Twid              // ~Left
	VaArg             // va_arg(Left, Type)
	Xor               // Left ^ Right
	XorEq             // Left ^= Right
	LCuBrk
	RCuBrk
)

var exprOpString = []string{
	Add:        "Add",
	AddEq:      "AddEq",
	Addr:       "Addr",
	And:        "And",
	AndAnd:     "AndAnd",
	AndEq:      "AndEq",
	Arrow:      "Arrow",
	Call:       "Call",
	Cast:       "Cast",
	CastInit:   "CastInit",
	Comma:      "Comma",
	Cond:       "Cond",
	Div:        "Div",
	DivEq:      "DivEq",
	Dot:        "Dot",
	Eq:         "Eq",
	EqEq:       "EqEq",
	Gt:         "Gt",
	GtEq:       "GtEq",
	Index:      "Index",
	Indir:      "Indir",
	Lsh:        "Lsh",
	LshEq:      "LshEq",
	Lt:         "Lt",
	LtEq:       "LtEq",
	Minus:      "Minus",
	Mod:        "Mod",
	ModEq:      "ModEq",
	Mul:        "Mul",
	MulEq:      "MulEq",
	Name:       "Name",
	Not:        "Not",
	NotEq:      "NotEq",
	Number:     "Number",
	Literal:    "Literal",
	Offsetof:   "Offsetof",
	Or:         "Or",
	OrEq:       "OrEq",
	OrOr:       "OrOr",
	Paren:      "Paren",
	Plus:       "Plus",
	PostDec:    "PostDec",
	PostInc:    "PostInc",
	PreDec:     "PreDec",
	PreInc:     "PreInc",
	Rsh:        "Rsh",
	RshEq:      "RshEq",
	SizeofExpr: "SizeofExpr",
	SizeofType: "SizeofType",
	String:     "String",
	Sub:        "Sub",
	SubEq:      "SubEq",
	Twid:       "Twid",
	VaArg:      "VaArg",
	Xor:        "Xor",
	XorEq:      "XorEq",
	LCuBrk:     "LCuBrk",
	RCuBrk:     "RCuBrk",
}

func (op ExprOp) String() string {
	if 0 <= int(op) && int(op) <= len(exprOpString) {
		return exprOpString[op]
	}
	return fmt.Sprintf("ExprOp(%d)", op)
}

// Prefix is an initializer prefix.
type Prefix struct {
	Span  Span
	Id    int
	Dot   Syntax // .Dot =
	XDecl *Decl  // for .Dot
	Index *Expr  // [Index] =
}

func (x *Prefix) GetId() int {
	return x.Id
}

func (x *Prefix) GetChildren() []Syntax {
	return []Syntax{}
}

func (x *Prefix) GetComments() *Comments {
	return nil
}

func (s *Prefix) GetSpan() Span {
	return Span{}
}

func (x *Prefix) String() string {
	if x.Dot.String() != "" {
		return "." + x.Dot.String()
	} else {
		return "[" + x.Index.String() + "]"
	}
}

// Init is an initializer expression.
type Init struct {
	SyntaxInfo
	Id     int
	Prefix []*Prefix // list of prefixes
	Expr   *Expr     // Expr
	Braced []*Init   // {Braced}

	XType *Type // derived type
}

func (x *Init) GetId() int {
	return x.Id
}

func (x *Init) GetChildren() []Syntax {
	lst := []Syntax{}
	for _, elem := range x.Prefix {
		lst = append(lst, elem)
	}
	if x.Expr != nil {
		lst = append(lst, x.Expr)
	}
	for _, elem := range x.Braced {
		lst = append(lst, elem)
	}
	return lst
}

func (x *Init) String() string {
	var p Printer
	p.hideComments = true
	p.printInit(x)
	return p.String()
}

// Walk traverses the syntax x, calling before and after on entry to and exit from
// each Syntax encountered during the traversal. In case of cross-linked input,
// the traversal never visits a given Syntax more than once.
func Walk(x Syntax, before, after func(Syntax)) {
	seen := map[Syntax]bool{
		nil:                     true,
		(*EmptyLiteral)(nil):    true,
		(*IntegerLiteral)(nil):  true,
		(*CharLiteral)(nil):     true,
		(*RealLiteral)(nil):     true,
		(*StringLiteral)(nil):   true,
		(*SymbolLiteral)(nil):   true,
		(*LanguageKeyword)(nil): true,
		(*Decl)(nil):            true,
		(*Init)(nil):            true,
		(*Type)(nil):            true,
		(*Expr)(nil):            true,
		(*Stmt)(nil):            true,
		(*Label)(nil):           true,
	}
	walk(x, before, after, seen)
}

func walk(x Syntax, before, after func(Syntax), seen map[Syntax]bool) {
	if x == nil {
		return
	}
	if seen[x] {
		return
	}
	seen[x] = true
	before(x)
	switch x := x.(type) {
	default:
		panic(fmt.Sprintf("walk: unexpected type %T", x))

	case *EmptyLiteral:
		//ok
	case *BooleanLiteral:
		//ok
	case *IntegerLiteral:
		//ok
	case *CharLiteral:
		//ok
	case *RealLiteral:
		//ok
	case *StringLiteral:
		//ok
	case *SymbolLiteral:
		//ok
	case *LanguageKeyword:
		//ok
	case *Prog:
		for _, d := range x.Decls {
			walk(d, before, after, seen)
		}

	case *Decl:
		walk(x.Type, before, after, seen)
		walk(x.Init, before, after, seen)
		walk(x.Body, before, after, seen)

	case *Init:
		for _, b := range x.Braced {
			walk(b, before, after, seen)
		}
		walk(x.Expr, before, after, seen)

	case *Type:
		walk(x.Base, before, after, seen)
		for _, d := range x.Decls {
			walk(d, before, after, seen)
		}
		walk(x.Width, before, after, seen)

	case *Expr:
		walk(x.Left, before, after, seen)
		walk(x.Text, before, after, seen)
		walk(x.Right, before, after, seen)
		for _, y := range x.LaunchParams {
			walk(y, before, after, seen)
		}
		for _, y := range x.Texts {
			walk(y, before, after, seen)
		}
		for _, y := range x.List {
			walk(y, before, after, seen)
		}
		walk(x.Type, before, after, seen)
		walk(x.Init, before, after, seen)
		for _, y := range x.Block {
			walk(y, before, after, seen)
		}

	case *Stmt:
		walk(x.Pre, before, after, seen)
		walk(x.Expr, before, after, seen)
		walk(x.Post, before, after, seen)
		walk(x.Decl, before, after, seen)
		walk(x.Body, before, after, seen)
		walk(x.Else, before, after, seen)
		walk(x.Text, before, after, seen)
		for _, y := range x.Block {
			walk(y, before, after, seen)
		}
		for _, y := range x.Labels {
			walk(y, before, after, seen)
		}

	case *Label:
		walk(x.Name, before, after, seen)
		walk(x.Expr, before, after, seen)
	}
	after(x)
}

// Preorder calls f for each piece of syntax of x in a preorder traversal.
func Preorder(x Syntax, f func(Syntax)) {
	Walk(x, f, func(Syntax) {})
}

// Preorder calls f for each piece of syntax of x in a postorder traversal.
func Postorder(x Syntax, f func(Syntax)) {
	Walk(x, func(Syntax) {}, f)
}
