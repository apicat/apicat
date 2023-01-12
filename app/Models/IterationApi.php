<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class IterationApi extends Model
{
    use HasFactory;

    protected $fillable = [
        'iteration_id', 'node_id', 'node_type'
    ];
}
