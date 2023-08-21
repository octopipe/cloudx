import React from "react"
import { Button, Form } from "react-bootstrap"
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/theme-github";

const WorkspaceSettings = () => {
  return (
    <div className="p-4" style={{width: '40rem'}}>
      <Form>
        <Form.Group className="mb-3">
          <Form.Label>Name</Form.Label>
          <Form.Control type="text" placeholder="Infra name..." />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Description</Form.Label>
          <Form.Control as="textarea" rows={3} />
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Override fields</Form.Label>
          <AceEditor
            mode="json"
            theme="github"
            onChange={() => console.log('changed')}
            name="UNIQUE_ID_OF_DIV"
            editorProps={{ $blockScrolling: true }}
            setOptions={{
              useWorker: false
            }}
            width="100%"
            height="200px"
          />
        </Form.Group>
      </Form>
      <div className="d-grid gap-2 mt-4">
        <Button variant="outline-danger">Delete workspace</Button>
      </div>
    </div>
  )
}

export default WorkspaceSettings