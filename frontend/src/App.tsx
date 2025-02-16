import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import Navbar from "./components/common/Navbar"
import About from "./pages/About"
//import Editor from "./pages/Editor"

const queryClient = new QueryClient()

function App() {
	return (
		<QueryClientProvider client={queryClient}> 
				<Navbar />
				{/*<Editor />*/}
				<About/>
		</QueryClientProvider>
	)
}

export default App
