<?php

namespace App\Exceptions;

use Exception;

class NotFoundException extends Exception
{
    /**
     * 未找到指定资源
     *
     * @return \Illuminate\Http\Response
     */
    public function render()
    {
        return response()->json(['status' => -404, 'msg' => 'Not Found']);
    }
}
