import './App.css';
import Editor from './Editor'
import { useState, useEffect, memo } from 'react';

const NodeList = memo(
  () => {
    const [nodeList, setnodeList] = useState([]);
    useEffect(() => {
      fetch("http://" + window.location.host + "/web/getNodeList")
        .then((res) => {
          return res.json();
        })
        .then((data) => {
          const nodes = data.split(',')
          // console.log(nodes)
          setnodeList(nodes)
        })
        .catch((error) => {
          console.log(error);
        });
    }, []);

    return (
      <div className="other">
        <h2>Node List</h2>
        <ul>
          {
            nodeList.map((nodeUrl, idx) => {
              return (
                <li>
                  <a
                    href={`http://${nodeUrl}`}
                    target="_blank"
                    rel="noreferrer"
                  >
                    Node #{idx}
                  </a>
                </li>
              )
            })

          }
        </ul>
      </div>
    )
  }
)

export default function App() {
  return (
    <div className="App">
      <h1>Meer</h1>
      <p>A collaborative peer-to-peer text editor</p>
      <Editor />
      <NodeList />
    </div>
  );
}

