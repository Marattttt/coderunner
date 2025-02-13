import { useState } from "react";
import ActionButton from "../components/common/ActionButton";
import CodeEditor from "../components/features/CodeEditor";
import LanguageSelect from "../components/features/LanguageSelect";
import {Languages }from "../constants"

const Editor = () => {
	let code = ''
	const [language, setLanguage] = useState(Languages[0])

	return (
		<div className="flex flex-col py-5 px-2 sm:p-8 gap-y-8 bg-bg-main w-full sm:h-screen">
			<div className="flex flex-wrap justify-between w-full">
				<LanguageSelect languages={Languages} onChange={(l) => setLanguage(l)}/>
				<ActionButton onClick={() => alert('action btn')} >
					Run!
				</ActionButton>
			</div>
			<div className="flex flex-col sm:flex-row gap-2 justify-items-stretch h-svh sm:h-full">
				<CodeEditor 
					code={code} 
					languageId={language.id}
					onChange={(c) => code = c}
					className="w-full h-full text-text-primary bg-bg-secondary"
				/>
				<div className="w-full h-full bg-bg-secondary"/>
			</div>
		</div >
	)
}

export default Editor;
