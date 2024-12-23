package parser

import (
	"fmt"
	"slices"
)

type bindingChain struct {
	binding Binding
	next    *bindingChain
}

func DelimitedParser(bindingParser Parser[Binding], separator rune) Parser[BindingList] {
	return Loop(nil,
		func(chain *bindingChain) Parser[Step[*bindingChain, BindingList]] {
			fmt.Printf("chain: %v\n", chain)
			if chain == nil {
				return Map(bindingParser,
					func(binding Binding) Step[*bindingChain, BindingList] {
						fmt.Printf("binding: %v\n", binding)
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
				fmt.Printf("binding: %v\n", b)

				return Step[*bindingChain, BindingList]{
					Accum: &bindingChain{binding: b, next: chain},
					Done:  false,
				}
			})

			var bindingSlice BindingList
			b := chain
			for {
				if b == nil {
					break
				}
				bindingSlice = slices.Insert(bindingSlice, 0, b.binding)
				b = b.next
			}
			return OneOf(
				extend,
				Succeed(Step[*bindingChain, BindingList]{Value: bindingSlice, Done: true}),
			)

		},
	)
}

func ConusmeBindingUntil(bindingParser Parser[Binding], terminator rune) Parser[BindingList] {
	return Loop(nil,
		func(chain *bindingChain) Parser[Step[*bindingChain, BindingList]] {
			fmt.Printf("chain: %v\n", chain)
			if chain == nil {
				return Map(bindingParser,
					func(binding Binding) Step[*bindingChain, BindingList] {
						fmt.Printf("binding: %v\n", binding)
						return Step[*bindingChain, BindingList]{Accum: &bindingChain{binding: binding}, Done: false}
					},
				)
			}

			s1 := StartSkipping(WhitespaceSkipParser)
			k1 := AppendKeeping(s1, bindingParser)
			s2 := AppendSkipping(k1, WhitespaceSkipParser)

			extend := Apply(s2, func(b Binding) Step[*bindingChain, BindingList] {
				fmt.Printf("binding: %v\n", b)

				return Step[*bindingChain, BindingList]{
					Accum: &bindingChain{binding: b, next: chain},
					Done:  false,
				}
			})

			var bindingSlice BindingList
			b := chain
			for {
				if b == nil {
					break
				}
				bindingSlice = slices.Insert(bindingSlice, 0, b.binding)
				b = b.next
			}
			return OneOf(
				extend,
				Succeed(Step[*bindingChain, BindingList]{Value: bindingSlice, Done: true}),
			)
		},
	)
}
