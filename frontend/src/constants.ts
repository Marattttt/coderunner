export interface Language {
	id: string
	name: string
}

export const Languages: Language[] = [
	{ id: 'go', name: "Go" },
	{ id: 'py', name: "Python" }
]

export const EditorUpdateDelay = 500;

export const CODERUNNER_API_URL = import.meta.env.VITE_CODERUNNER_API_URL
