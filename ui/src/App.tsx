import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Main from './modules/Main'
import SharedInfras from './modules/SharedInfras'
import ConnectionInterfaces from './modules/CloudAccounts'
import CloudAccounts from './modules/ConnectionInterfaces'
import SharedInfraView from './modules/SharedInfrasView'

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Main />}>
          <Route path='' element={<Navigate to="/shared-infras" replace />} />
          <Route path='/shared-infras' element={<SharedInfras />} />
          <Route path='/shared-infras/:name' element={<SharedInfraView />}/>
          <Route path='/connection-interfaces' element={<ConnectionInterfaces />}/>
          <Route path='/cloud-accounts' element={<CloudAccounts />}/>
        </Route>
        <Route path="login" element={<div>Login</div>} />
        <Route path="*" element={<div>Not Found</div>} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
