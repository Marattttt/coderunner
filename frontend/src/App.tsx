import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { Route, Routes } from "react-router-dom"
import Navbar from "./components/common/Navbar"
import About from "./pages/About/About"
import Editor from "./pages/Editor/Editor"
import Footer from "./components/common/Footer"

const queryClient = new QueryClient()

function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<div className="flex flex-col min-h-screen">
				<Navbar />
					<main className="flex-grow">
						<Routes>
							<Route path="/" element={<Editor />} />
							<Route path="/about" element={<About />} />
						</Routes>
					</main>
				<Footer />
			</div>
		</QueryClientProvider>
	)
}

export default App
