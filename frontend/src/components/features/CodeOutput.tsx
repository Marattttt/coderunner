import React from "react"

interface CodeOutputParams {
	stdout?: string
	stderr?: string
	error?: string
	className?: string
}

const CodeOutput: React.FC<CodeOutputParams> = ({stdout, stderr, error, className}) => {
	return (<div
		className={className + " flex flex-col overflow-auto gap-y-4"}>
		{stdout && 
			<pre className="text-text-stdout">
				------<br/>
				{stdout}
			</pre>
		}

		{stderr && 
			<pre className="text-text-stderr">
				------<br/>
				{stderr}
			</pre>
		}
		
		{error && 
			<pre className="text-text-error">
				Error:<br/>
				{error}
			</pre>
		}
	</div>)
}

export default CodeOutput
