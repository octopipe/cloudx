import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Main from './modules/Main'
import SharedInfras from './modules/SharedInfras'
import ConnectionInterfaces from './modules/CloudAccounts'
import CloudAccounts from './modules/ConnectionInterfaces'

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Main />}>
          <Route path='/' element={<SharedInfras />}/>
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
