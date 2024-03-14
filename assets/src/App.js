import './App.css';
import Editor from './Editor'
import { useState } from 'react';
import {QueryClient, QueryClientProvider} from "react-query";

function App() {
  const queryClient = new QueryClient();

  return (
  <QueryClientProvider client={queryClient}>
    <div className="App">
      
      <Editor className="editor"/>   
    </div>
    </QueryClientProvider>
  );
}

export default App;
