import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { memo } from 'react';
import { Alert, Modal, Toast } from 'react-bootstrap';
import { Handle, Position } from 'reactflow';


export default memo(({ data, selected }: any) => {
  return (
    <>
      <Handle type="target" position={Position.Left} />
      <div className='d-flex'>
        <div className='bar'></div>
        <div className='d-flex flex-column content'>
          <div>
            <strong>{data.label}</strong>
          </div>
          <div><FontAwesomeIcon icon="clock" />{' '}30s</div>
          {data?.status.indexOf("ERROR") !== -1 && (
            <a 
              className='link-danger'
              style={{ cursor: 'pointer' }}
            >
              <FontAwesomeIcon icon="circle-exclamation" color='red' /> 1 Error found
            </a>
          )}
        </div>
      </div>
      <Handle type="source" position={Position.Right} />
      
    </>
  );
});