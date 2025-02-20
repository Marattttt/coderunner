import React, { useRef, useState } from "react"
import Prism from 'prismjs'
import Editor from "react-simple-code-editor";
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-go'

interface editorProps {
	code: string,
	languageId: string,
	className?: string,
	onChange?: (code: string) => void
}

const prismLanguages: Record<string, string> = {
	'go': 'go',
	'py': 'python',
}

const CodeEditor: React.FC<editorProps> = ({
	code,
	languageId,
	className,
	onChange,
}) => {
	const [rawCode, setCode] = useState(code)

	const containerRef = useRef(null);

	const handleContainerClick = (_: any) => {
		// Editor component is a thin wrapper around textarea,
		// so even if this is a hack, it works well
		//@ts-ignore
		const textarea = containerRef.current!.querySelector('textarea');
		if (textarea) {
			textarea.focus();
		}
	};

	const prismLang = prismLanguages[languageId]

	const addHighlght = (code: string) => {
		return Prism.highlight(
			code,
			Prism.languages[prismLang],
			prismLang
		);
	}

	return (
		<div
			className={(className ?? '') + ` font-code overflow-y-auto language-${prismLang}`}
			ref={containerRef}
			onClick={handleContainerClick}
		>
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
