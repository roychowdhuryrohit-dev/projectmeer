import './Editor.css'
import {createEditor, $createTextNode, $getRoot, $getSelection, $createParagraphNode} from 'lexical';
import {useDispatch, useSelector} from 'react-redux';
import {useParams} from 'react-router-dom';
import {useEffect, useState, useRef} from 'react';
import axios from 'axios';
import {LexicalComposer} from '@lexical/react/LexicalComposer';
import {PlainTextPlugin} from '@lexical/react/LexicalPlainTextPlugin';
import {ContentEditable} from '@lexical/react/LexicalContentEditable';
import {HistoryPlugin} from '@lexical/react/LexicalHistoryPlugin';
import {OnChangePlugin} from '@lexical/react/LexicalOnChangePlugin';
import {useLexicalComposerContext} from '@lexical/react/LexicalComposerContext';
import LexicalErrorBoundary from '@lexical/react/LexicalErrorBoundary';
const JsDiff = require('diff');


function MyOnChangePlugin({ onChange }) {
  const [editor] = useLexicalComposerContext();

  useEffect(() => {
    return editor.registerUpdateListener(({ editorState }) => {
      onChange(editorState);
    });
  }, [editor, onChange]);

  return null;
}

const initialConfig = {
  namespace: 'MyEditor',
  onError: (error)=>console.error(error),

};



function Editor() {
  const [editorState, setEditorState] = useState(null)
  const [previousEditorState, setPreviousEditorState] = useState(null);
  const updateEditorState = (editorState)=> {
  

    editorState.read(() => {
    
      if (previousEditorState===null || previousEditorState === undefined) {
        setPreviousEditorState($getRoot().getTextContent());
        return;
      }
      
      const oldText = previousEditorState;
      const newText = $getRoot().getTextContent();
      const diff = JsDiff.diffChars(oldText, newText);
      
      function getChangeStartIndex (change) {
        let index = 0;
      
        for (let i = 0; i < diff.length; i++) {
            if (diff[i] === change) {
            return index;
            }   
      
            if (!diff[i].removed) {
                index += diff[i].value.length;
            }
        }
      
        return -1; // Change not found
      }
      
      for (let d of diff) {
      if (d.added!==undefined || d.removed!==undefined) {
        console.log(d);
        const startIndex = getChangeStartIndex(d);
        console.log("Start index in old text:", startIndex);
        if (d.added) {
          //call insertextapi
          const fetchData = async () => {
            try {
              const jsonData = {
                index: startIndex,
                text: d.value
              };

              const requestOptions = {
                method: 'GET', 
                headers: {
                  'Content-Type': 'application/json',
                },
                body: JSON.stringify(jsonData),
                };
               
               const response = await fetch(`/web/insertText`, requestOptions);
               const insertText = await response.json();
               setEditorState(insertText.data);
            } catch (error) {
              console.error(error);
            }
          }
          fetchData();
          
        } else {
          //call removetextapi
          const jsonData = {
            index: startIndex,
            count: d.count
          };

          const requestOptions = {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(jsonData),
            };

            try {
              const response = fetch('/web/removeText', requestOptions);
              const deletedText = response.json();

              setEditorState(deletedText.json());
            } catch(error) {
              console.error(error);
            }
            }
          }
        }
      })
//After this, get new text from GET request and update state
      
      // newText = fetch new text from get api

      // const root = editorState.$getRoot();
      // root.clear();
      // const paragraph = $createParagraphNode();
      // paragraph.append($createTextNode(newText));
      // root.append(paragraph);
      
      
      setPreviousEditorState(newText);
    });
  }

  return (
   <LexicalComposer initialConfig={initialConfig} >
    <div className='editor-inner'>
      <PlainTextPlugin
        contentEditable={<ContentEditable />}
        placeholder={<div>Enter some text...</div>}
        ErrorBoundary={LexicalErrorBoundary}
      />
      <HistoryPlugin />
      <OnChangePlugin onChange={updateEditorState}/>
    </div>
    </LexicalComposer>
  );
}

export default Editor;