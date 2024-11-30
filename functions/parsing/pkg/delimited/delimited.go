package delimited

import (
	. "dave.internal/pkg/parser"
)

type BindingList struct {
	binding Binding
	next    *BindingList
}

func DelimitedParser(bindingParser Parser[Binding], separator string) Parser[[]Binding] {
	return Loop(nil,
		func(bindings *BindingList) Parser[Step[*BindingList, []Binding]] {
			if bindings == nil {
				return Map(bindingParser,
					func(binding Binding) Step[*BindingList, []Binding] {
						return Step[*BindingList, []Binding]{Accum: &BindingList{binding: binding}, Done: false}
					},
				)
			}
			s := StartSkipping(WhitespaceSkipParser)
			s2 := AppendKeeping(s, bindingParser)
			extend := Apply(s2, func(b Binding) Step[*BindingList, []Binding] {
				return Step[*BindingList, []Binding]{
					Accum: &BindingList{binding: b, next: bindings},
					Done:  false,
				}
			})

			var bindingSlice []Binding
			b := bindings
			for {
				if b == nil {
					break
				}
				bindingSlice = append(bindingSlice, b.binding)
				b = b.next
			}
			return OneOf(
				extend,
				Succeed(Step[*BindingList, []Binding]{Value: bindingSlice, Done: true}),
			)

		},
	)
}
