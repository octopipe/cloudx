import { faBorderTopLeft } from '@fortawesome/free-solid-svg-icons';
import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'

export default memo(({ data, isConnectable }: any) => {
  return (
    <>
      <div className=''>
        <strong>{data.label}</strong><br/>
        <FontAwesomeIcon className='mt-2' icon="diagram-project" size='2x' />
      </div>
      <Handle
        type="source"
        position={Position.Right}
        id="a"
        style={{ background: '#555' }}
        isConnectable={isConnectable}
      />
    </>
  );
});