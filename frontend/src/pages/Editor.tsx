import { useState } from "react";
import ActionButton from "../components/common/ActionButton";
import CodeEditor from "../components/features/CodeEditor";
import LanguageSelect from "../components/features/LanguageSelect";
import {Languages }from "../constants"

const Editor = () => {
	let code = ''
	const [language, setLanguage] = useState(Languages[0])

	return (
		<div className="flex flex-col p-10 sm:p-10 gap-y-32 bg-bg-main w-full">
			<div className="flex flex-wrap justify-between w-full">
				<LanguageSelect languages={Languages} onChange={(l) => setLanguage(l)}/>
				<ActionButton onClick={() => alert('action btn')} >
					Run!
				</ActionButton>
			</div>
			<div className="flex flex-col sm:flex-row gap-4 sm:gap-8 justify-items-stretch">
				<CodeEditor 
					code={code} 
					languageId={language.id}
					onChange={(c) => code = c}
					className="w-full h-24 text-text-primary bg-bg-secondary"
				/>
				<div className="h-24 w-full bg-bg-secondary" />
			</div>
		</div >
	)
}

export default Editor;
