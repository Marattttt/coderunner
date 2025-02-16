import { useState } from "react";
import ActionButton from "../components/common/ActionButton";
import CodeEditor from "../components/features/CodeEditor";
import LanguageSelect from "../components/features/LanguageSelect";
import { Languages } from "../constants"
import { useQuery } from "@tanstack/react-query";
import runCode from "../api/coderunner/runCode";
import CodeOutput from "../components/features/CodeOutput";

const Editor = () => {
	// Not a useState, due to CodeEditor handling the rerender on its own
	const [code, setCode] = useState('')
	const [language, setLanguage] = useState(Languages[0])

	const outputQuery = useQuery({
		queryKey: [],
		queryFn: () => runCode(code, language.id),
		enabled: false,
		retry: false,
	})

	const executeCode = () => {
		outputQuery.refetch({
			cancelRefetch: true
		})
	}

	return (
		<div className="flex flex-col py-8 px-2 sm:px-8 gap-y-4 bg-bg-main w-full sm:h-screen">
			<div className="flex flex-wrap justify-between w-full">
				<LanguageSelect languages={Languages} onChange={(l) => setLanguage(l)} />
				<ActionButton
					onClick={executeCode}
					disabled={outputQuery.isLoading}
				>
					Run!
				</ActionButton>
			</div>
			<div className="flex flex-col sm:flex-row gap-2 justify-items-stretch h-svh sm:h-full">
				<CodeEditor
					code={code}
					languageId={language.id}
					onChange={(c) => setCode(c)}
					className="w-full h-full py-4 px-2 text-text-primary bg-bg-secondary"
				/>
				<CodeOutput
					stdout={outputQuery.data?.stdout}
					stderr={outputQuery.data?.stderr}
					error={outputQuery.error?.message}
					className="h-full w-full py-4 px-2 bg-bg-secondary"
				/>
			</div>
		</div >
	)
}

export default Editor;
