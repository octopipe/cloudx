import React from 'react'
import './index.css'
import { Card, ListGroup, Nav } from 'react-bootstrap'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

const AddPanel = ({ node, onClose }: any) => {
  const onDragStart = (event: any, nodeType: any) => {
    event.dataTransfer.setData('application/reactflow', nodeType);
    event.dataTransfer.effectAllowed = 'move';
  };

  return (
    <Card className='shared-infra-diagram__add-panel'>
      <Card.Body>
        <Nav>
          <Nav.Item onDragStart={(event: any) => onDragStart(event, 'default')} draggable><FontAwesomeIcon icon="box" /></Nav.Item>
          <Nav.Item onDragStart={(event: any) => onDragStart(event, 'aws')} draggable><FontAwesomeIcon icon={["fab", "aws"]} /></Nav.Item>
          <Nav.Item  onDragStart={(event: any) => onDragStart(event, 'connection-interface')} draggable><FontAwesomeIcon icon="diagram-project"/></Nav.Item>
        </Nav>
      </Card.Body>
    </Card>
  )
  
}

export default AddPanel