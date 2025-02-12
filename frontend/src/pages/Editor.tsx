const Editor = () => {
	return (
		<div className="flex flex-col p-10 sm:p-10 gap-y-32 bg-bg-main w-full">
			<div className="flex flex-wrap justify-between w-full">
				<span> Language select </span>
				<span> Run button </span>
			</div>
			<div className="flex flex-col sm:flex-row gap-4 sm:gap-10 justify-items-stretch">
				<div className="h-24 w-full bg-bg-secondary" />
				<div className="h-24 w-full bg-bg-secondary" />
			</div>
		</div>
	)
}

export default Editor;
