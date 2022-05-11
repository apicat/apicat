<?php

namespace App\Exceptions\Mock;

use Exception;

/**
 * 404
 */
class NotFoundException extends Exception
{
    public function render()
    {
        return response()->json([
            'status' => -1,
            'msg' => 'Not Found'
        ], 404);
    }
}
