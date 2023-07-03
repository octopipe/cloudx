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
    <div className='shared-infra-diagram__add-panel'>
        <div style={{cursor: "grab"}} onDragStart={(event: any) => onDragStart(event, 'default')} draggable><FontAwesomeIcon icon="box" /></div>
        <div style={{cursor: "grab"}} onDragStart={(event: any) => onDragStart(event, 'aws')} draggable><FontAwesomeIcon icon={["fab", "aws"]} /></div>
        <div style={{cursor: "grab"}} onDragStart={(event: any) => onDragStart(event, 'connection-interface')} draggable><FontAwesomeIcon icon="diagram-project"/></div>
    </div>
  )
  
}

export default AddPanel