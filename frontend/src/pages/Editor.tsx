import ActionButton from "../components/common/ActionButton";
import LanguageSelect from "../components/features/LanguageSelect";
import {Languages }from "../constants"

const Editor = () => {
	return (
		<div className="flex flex-col p-10 sm:p-10 gap-y-32 bg-bg-main w-full">
			<div className="flex flex-wrap justify-between w-full">
				<LanguageSelect languages={Languages} onChange={(s) => console.log(`${s} selected`)}/>
				<ActionButton onClick={() => alert('action btn')} >
					Run!
				</ActionButton>
			</div>
			<div className="flex flex-col sm:flex-row gap-4 sm:gap-8 justify-items-stretch">
				<div className="h-24 w-full bg-bg-secondary" />
				<div className="h-24 w-full bg-bg-secondary" />
			</div>
		</div >
	)
}

export default Editor;
