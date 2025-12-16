import { Link } from '@tanstack/react-router'

import { useState } from 'react'
import {
  Cookie
} from 'lucide-react'

export default function Header() {

  return (
    <>
      <header className="p-4 flex items-center bg-gray-800 text-white shadow-lg">
        <h1 className="ml-4 text-xl font-semibold flex flex-row items-center">
          <Cookie />
          <p>JAFA</p>
        </h1>
        <nav>
        </nav>
      </header>
    </>
  )
}
