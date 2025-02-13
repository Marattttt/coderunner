export interface Language {
	id: string
	name: string
}

export const Languages: Language[] = [
	{ id: 'go', name: "Go" },
	{ id: 'js', name: "JavaScript" }
]

export const EditorUpdateDelay = 500;
