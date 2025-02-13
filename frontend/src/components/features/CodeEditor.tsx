import React, { useState } from "react"
import Prism from 'prismjs'
import Editor from "react-simple-code-editor";
import 'prismjs/components/prism-clike'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-go'

interface editorProps {
	code: string,
	languageId: string,
	className?: string,
	onChange?: (code: string) => void
}

const prismLanguages: Record<string, string> = {
	'go': 'go',
	'js': 'js',
}

const CodeEditor: React.FC<editorProps> = ({
	code,
	languageId,
	className,
	onChange,
}) => {
	const [rawCode, setCode] = useState(code)

	const addHighlght = (code: string) => {
		return Prism.highlight(
			code,
			Prism.languages[languageId],
			languageId
		);
	}

	return (
		<div className={className ?? '' + ` language-${prismLanguages[languageId]}`}>
			<Editor
				value={rawCode}
				onValueChange={(c) => {
					setCode(c);
					onChange && onChange(c)
				}}
				highlight={addHighlght}
			/>
		</div>
	)

}


export default CodeEditor
