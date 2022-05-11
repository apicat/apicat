<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Support\Facades\Auth;

class LogoutController extends Controller
{
    /**
     * Log the user out of the application.
     *
     * @return array
     */
    public function index()
    {
        Auth::guard('api')->logout();
        return ['status' => 0, 'msg' => '退出成功'];
    }
}
