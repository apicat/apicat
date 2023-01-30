<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\ApiDoc\ApiDocHistoryRepository;
use App\Repositories\User\UserRepository;
use Illuminate\Http\Request;

class ApiDocHistoryController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    public function histories(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1'],
        ]);

        $records = ApiDocHistoryRepository::list($request->input('doc_id'));
        if ($records->isEmpty()) {
            return [
                'status' => 0,
                'msg' => '',
                'data' => []
            ];
        }

        $users = UserRepository::idNameArr(true);

        $result = [];
        foreach ($records as $record) {
            $month = $record->last_updated_at->format('Y-m');
            if (!isset($result[$month])) {
                $result[$month] = [
                    [
                        'id' => $record->id,
                        'title' => $record->last_updated_at->format('m月d日 H:i') . '(' . (isset($users[$record->last_user_id]) ? $users[$record->last_user_id] : '') . ')',
                        'type' => 1
                    ]
                ];
            } else {
                $result[$month][] = [
                    'id' => $record->id,
                    'title' => $record->last_updated_at->format('m月d日 H:i') . '(' . (isset($users[$record->last_user_id]) ? $users[$record->last_user_id] : '') . ')',
                    'type' => 1
                ];
            }
        }

        $result2 = [];
        foreach ($result as $k => $v) {
            $result2[] = [
                'id' => 0,
                'title' => str_replace('-', '年', $k) . '月',
                'type' => 0,
                'sub_nodes' => $v
            ];
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => $result2
        ];
    }
}
