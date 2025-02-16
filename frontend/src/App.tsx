import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { Navigate, Route, Routes } from "react-router-dom"
import Navbar from "./components/common/Navbar"
import About from "./pages/About"
import Editor from "./pages/Editor"

const queryClient = new QueryClient()

function App() {
	return (
		<QueryClientProvider client={queryClient}> 
			<Navbar />
			<Routes>
				<Route path="/" element={ <Navigate to="/editor" /> }/>
				<Route path="/editor" element={ <Editor /> } />
				<Route path="/about" element={ <About /> } />
			</Routes>
		</QueryClientProvider>
	)
}

export default App
