package parser

type bindingChain struct {
	binding Binding
	next    *bindingChain
}

func DelimitedParser(bindingParser Parser[Binding], separator rune) Parser[BindingList] {
	return Loop(nil,
		func(bindings *bindingChain) Parser[Step[*bindingChain, BindingList]] {
			if bindings == nil {
				return Map(bindingParser,
					func(binding Binding) Step[*bindingChain, BindingList] {
						return Step[*bindingChain, BindingList]{Accum: &bindingChain{binding: binding}, Done: false}
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

			extend := Apply(s3, func(b Binding) Step[*bindingChain, BindingList] {
				return Step[*bindingChain, BindingList]{
					Accum: &bindingChain{binding: b, next: bindings},
					Done:  false,
				}
			})

			var bindingSlice BindingList
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
				Succeed(Step[*bindingChain, BindingList]{Value: bindingSlice, Done: true}),
			)

		},
	)
}
