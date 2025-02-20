import React from "react"

interface ContactProps {
	className?: string
	logo: React.ReactNode
	text: string
	url: string
}

const Contact: React.FC<ContactProps> = ({logo, text, url, className}) => {
	return (<a href={url} target="_blank">
		<div className={'flex flex-row justify-between items-center gap-4 text-xl font-bold ' + className}>
			{logo}
			{text}
		</div>
	</a>)
}

export default Contact
