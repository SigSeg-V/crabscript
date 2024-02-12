package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "Global"
	LocalScope   SymbolScope = "Local"
	BuiltinScope SymbolScope = "Builtin"
	FreeScope    SymbolScope = "Free"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
	FreeSymbols    []Symbol
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s, FreeSymbols: nil}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: s.numDefinitions,
	}

	if s.Outer != nil {
		// we are in global (most outer scope)
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		if !ok {
			return obj, ok
		}

		if obj.Scope == GlobalScope || obj.Scope == BuiltinScope {
			return obj, ok
		}

		free := s.DefineFree(obj)
		return free, true
	}
	return obj, ok
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	sym := Symbol{
		Name:  name,
		Index: index,
		Scope: BuiltinScope,
	}

	s.store[name] = sym
	return sym
}

func (s *SymbolTable) DefineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	sym := Symbol{
		Name:  original.Name,
		Scope: FreeScope,
		Index: len(s.FreeSymbols) - 1,
	}

	s.store[original.Name] = sym
	return sym
}
