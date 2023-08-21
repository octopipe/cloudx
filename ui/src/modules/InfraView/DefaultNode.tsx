import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';


export default memo(({ data, selected }: any) => {
  return (
    <>
      <Handle type="target" position={Position.Left} />
      <div className='d-flex'>
        <div className='content'>
          <strong>{data.label}</strong>
        </div>
       
      </div>
      <Handle type="source" position={Position.Right} />
    </>
  );
});