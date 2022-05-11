<?php

namespace App\Http\Controllers\Mock;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use App\Models\MockPath;
use App\Models\ApiDoc;
use App\Modules\Mock\Parser\ContentParser;
use App\Exceptions\Mock\NotFoundException;

class MockController extends Controller
{
    public function index(Request $request, $projectId = '', $path = '')
    {
        if (!$projectId or !is_numeric($projectId)) {
            throw new NotFoundException;
        }

        if (!$path) {
            throw new NotFoundException;
        }

        $path = '/' . $path;

        $record = MockPath::where([
            ['project_id', $projectId],
            ['path', $path]
        ])->first();

        if (!$record) {
            throw new NotFoundException;
        }

        if (!$doc = ApiDoc::find($record->doc_id)) {
            throw new NotFoundException;
        }

        $header = [];
        if ($record->format == 'json') {
            $body = '{}';
        } else {
            $body = '<?xml version="1.0" encoding="UTF-8" ?>';
        }
        
        if (!$doc->content) {
            return $this->resp($header, $body, $record->format);
        }

        $data = json_decode($doc->content);
        if (!isset($data->content) or !is_array($data->content)) {
            return $this->resp($header, $body, $record->format);
        }

        foreach ($data->content as $content) {
            if (!isset($content->type)) {
                continue;
            }

            if ($content->type == 'http_api_response_parameter') {
                $result = ContentParser::mock($content->attrs, $record->format);
                return $this->resp($result['mock_header'], $result['mock_body'], $record->format);
            }
        }
        
        return $this->resp($header, $body, $record->format);
    }

    protected function resp($header, $body, $format)
    {
        $response = response($body);

        if ($header) {
            foreach ($header as $k => $v) {
                $response->header($k, $v);
            }
        }

        if ($format == 'json') {
            $response->header('Content-Type', 'application/json');
        } else {
            $response->header('Content-Type', 'text/xml');
        }

        return $response;
    }
}
