import React, { useRef, useState } from "react"
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
	const addHighlght = (code: string) => {
		return Prism.highlight(
			code,
			Prism.languages[languageId],
			languageId
		);
	}

	return (
		<div 
			className={className ?? '' + ` language-${prismLanguages[languageId]}`}
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
