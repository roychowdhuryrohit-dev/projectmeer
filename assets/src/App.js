import './App.css';
import Editor from './Editor'

export default function App() {
  return (
    <div className="App">
      <h1>Meer</h1>
      <p>A collaborative peer-to-peer text editor</p>
      <Editor />
      <div className="other">
        <h2>Node List</h2>
        <ul>
          <li>
            <a
              href="https://codesandbox.io/s/lexical-rich-text-example-5tncvy"
              target="_blank"
              rel="noreferrer"
            >
              Rich text example
            </a>
          </li>
        </ul>
      </div>
    </div>
  );
}

