import EllipseSvg from '../../assets/Ellipse.svg?react'
import { Link } from "react-router-dom"

const Navbar = () => {
	return (
		<nav className="flex justify-between w-full px-4 sm:px-10 h-[64px] top-0 left-0 font-heading text-text-primary bg-bg-dark">
			<Link
				to="/"
				className="bg-none my-auto border-none font-bold text-2xl"
			>
				Coderunner
			</Link>

			<div className="flex gap-8 align-middle justify-center bg-none border-none font-bold text-lg" >
				<Link to="/editor" className="m-auto">
					Editor
				</Link>
				<EllipseSvg className="m-auto" />
				<Link to="/about" className="m-auto">
					About
				</Link>
			</div>
		</nav>
	)
}

export default Navbar
