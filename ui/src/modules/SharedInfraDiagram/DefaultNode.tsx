import { faBorderTopLeft } from '@fortawesome/free-solid-svg-icons';
import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';

const getDuration = (startedAt: any, finishedAt: any) => {
  const d1: any = new Date(startedAt)
  const d2: any = new Date(finishedAt)
  const diff = d2 - d1

  if (diff > 60e3) 
    return `${Math.floor(diff / 60e3)} minutes`

  return `${Math.floor(diff / 1e3)} seconds`
}

const colorByStatus: any = {
  'terraform': '#7b00db',
}

export default memo(({ data, isConnectable }: any) => {
  return (
    <>
      <Handle
        type="target"
        position={Position.Left}
        style={{ background: '#555' }}
        onConnect={(params) => console.log('handle onConnect', params)}
        isConnectable={isConnectable}
      />
      <div>
        <div
          style={{
            background: colorByStatus[data?.type],
            color: '#fff',
            padding: '5px',
            borderTopRightRadius: '5px',
            borderTopLeftRadius: '5px'

          }}
        >
          {data.label}
        </div>
        { data?.startedAt && data?.finishedAt && (
          <div style={{padding: '10px'}}>

            {getDuration(data?.startedAt, data?.finishedAt)}
            {/* {(data?.status === "ERROR" || data?.status === "FAILED") && (
              <Button variant='danger'>See error</Button>
            )} */}

          </div>
        )}
        
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