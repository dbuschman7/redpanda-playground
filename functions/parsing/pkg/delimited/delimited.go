package delimited

import (
	. "dave.internal/pkg/parser"
)

type BindingList struct {
	binding Binding
	next    *BindingList
}

func DelimitedParser(bindingParser Parser[Binding], separator rune) Parser[[]Binding] {
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
			s1 := AppendSkipping(s, ConsumeIf(func(r rune) bool {
				return r == rune(separator)
			}))
			s2 := AppendSkipping(s1, WhitespaceSkipParser)
			k1 := AppendKeeping(s2, bindingParser)
			s3 := AppendSkipping(k1, WhitespaceSkipParser)

			extend := Apply(s3, func(b Binding) Step[*BindingList, []Binding] {
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
