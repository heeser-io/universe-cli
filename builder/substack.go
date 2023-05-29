package builder

func (b *Builder) GetSubBuilder() []*Builder {
	builderArr := []*Builder{}
	for _, substack := range b.stack.Stacks {
		subBuilder, err := New(substack.BasePath, false)
		if err != nil {
			panic(err)
		}
		builderArr = append(builderArr, subBuilder)
	}

	return builderArr
}

// func (b *Builder) buildSubstack() {
// 	for _, substack := range b.stack.Stacks {
// 		subBuilder, err := New(substack.BasePath)
// 		if err != nil {
// 			panic(err)
// 		}

// 		BuildStack(subBuilder)
// 	}
// }
