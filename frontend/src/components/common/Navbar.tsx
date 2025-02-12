import React from "react"
import EllipseSvg from '../../assets/Ellipse.svg?react'

interface linkbtnprops {
	text: string
}

const LinkButton: React.FC<linkbtnprops> = ({ text }) => {
	return (
		<button
			onClick={() => console.log(`${text} pressed`)}
		>
			{text}
		</button>
	)
}
const Navbar = () => {
	return (
		<nav className="flex justify-between w-full px-10 h-[64px] top-0 left-0 text-text-primary bg-bg-dark">
			<button
				className="bg-none border-none font-bold text-2xl"
				onClick={() => console.log('Main button pressed')}>
				Coderunner
			</button>

			<div
				className="flex gap-8 align-middle justify-center bg-none border-none font-bold text-lg"
			>
				<LinkButton text="Editor" />
				<EllipseSvg className="m-auto" />
				<LinkButton text="About" />
			</div>

		</nav>
	)
}

export default Navbar
