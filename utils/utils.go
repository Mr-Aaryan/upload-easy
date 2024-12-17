package utils

var FolderStack []string

func PushToStack(parentId string) {
	FolderStack = append(FolderStack, parentId)
}

func PopFromStack() string {
	if len(FolderStack) == 0 {
		return ""
	}
	last := FolderStack[len(FolderStack)-1]
	FolderStack = FolderStack[:len(FolderStack)-1] //slice(0:last element)--returns upto 2nd last element and removes last one
	return last
}

func PeekStack() string {
	if len(FolderStack) == 0 {
		return ""
	}
	return FolderStack[len(FolderStack)-1]
}