import React from "react"

interface CodeOutputParams {
	stdout?: string
	stderr?: string
	error?: string
}

const CodeOutput: React.FC<CodeOutputParams> = ({stdout, stderr, error}) => {
	return (<div className="h-full w-full flex flex-col overflow-auto gap-y-8 p-8 bg-bg-secondary">
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
